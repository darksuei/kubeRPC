package webhook

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const (
	annotationEnabled = "kuberpc.suei.io/enabled"
	annotationService = "kuberpc.suei.io/service"
	annotationPort    = "kuberpc.suei.io/port"
	annotationHost    = "kuberpc.suei.io/host"
	defaultRPCPort    = "7749"
)

type admissionReview struct {
	APIVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Request    *admissionRequest  `json:"request,omitempty"`
	Response   *admissionResponse `json:"response,omitempty"`
}

type admissionRequest struct {
	UID       string          `json:"uid"`
	Namespace string          `json:"namespace"`
	Object    json.RawMessage `json:"object"`
}

type admissionResponse struct {
	UID       string  `json:"uid"`
	Allowed   bool    `json:"allowed"`
	PatchType *string `json:"patchType,omitempty"`
	Patch     []byte  `json:"patch,omitempty"`
}

type podMeta struct {
	Metadata struct {
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Containers []struct {
			Name string   `json:"name"`
			Env  []envVar `json:"env"`
		} `json:"containers"`
	} `json:"spec"`
}

type envVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type jsonPatch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

func Mutate(w http.ResponseWriter, r *http.Request) {
	coreURL := os.Getenv("KUBERPC_CORE_URL")
	if coreURL == "" {
		slog.Error("webhook: KUBERPC_CORE_URL not set")
		http.Error(w, "webhook misconfigured", http.StatusInternalServerError)
		return
	}

	var review admissionReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		slog.Error("webhook: decode failed", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if review.Request == nil {
		http.Error(w, "missing admission request", http.StatusBadRequest)
		return
	}

	review.Response = mutate(review.Request, coreURL)
	review.Request = nil

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(review); err != nil {
		slog.Error("webhook: encode response failed", "error", err)
	}
}

func mutate(req *admissionRequest, coreURL string) *admissionResponse {
	allow := &admissionResponse{UID: req.UID, Allowed: true}

	var pod podMeta
	if err := json.Unmarshal(req.Object, &pod); err != nil {
		slog.Error("webhook: unmarshal pod failed", "error", err)
		return allow
	}

	ann := pod.Metadata.Annotations
	if ann[annotationEnabled] != "true" {
		return allow
	}

	slog.Info("webhook: mutating pod", "namespace", req.Namespace)

	toInject := buildEnvVars(ann, req.Namespace, coreURL)
	if len(toInject) == 0 {
		return allow
	}

	patches := buildPatches(pod, toInject)
	if len(patches) == 0 {
		return allow
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		slog.Error("webhook: marshal patch failed", "error", err)
		return allow
	}

	pt := "JSONPatch"
	return &admissionResponse{
		UID:       req.UID,
		Allowed:   true,
		PatchType: &pt,
		Patch:     patchBytes,
	}
}

func buildEnvVars(ann map[string]string, namespace, coreURL string) []envVar {
	vars := []envVar{
		{Name: "KUBERPC_CORE_URL", Value: coreURL},
	}

	svcName, ok := ann[annotationService]
	if !ok || svcName == "" {
		return vars
	}

	host := fmt.Sprintf("%s.%s.svc.cluster.local", svcName, namespace)
	if custom := ann[annotationHost]; custom != "" {
		host = custom
	}

	port := defaultRPCPort
	if p := ann[annotationPort]; p != "" {
		port = p
	}

	return append(vars,
		envVar{Name: "KUBERPC_SERVICE_NAME", Value: svcName},
		envVar{Name: "KUBERPC_HOST", Value: host},
		envVar{Name: "KUBERPC_PORT", Value: port},
	)
}

func buildPatches(pod podMeta, toInject []envVar) []jsonPatch {
	var patches []jsonPatch

	for i, container := range pod.Spec.Containers {
		existing := make(map[string]bool, len(container.Env))
		for _, e := range container.Env {
			existing[e.Name] = true
		}

		var fresh []envVar
		for _, e := range toInject {
			if !existing[e.Name] {
				fresh = append(fresh, e)
			}
		}
		if len(fresh) == 0 {
			continue
		}

		base := fmt.Sprintf("/spec/containers/%d/env", i)

		if len(container.Env) == 0 {
			patches = append(patches, jsonPatch{Op: "add", Path: base, Value: fresh})
		} else {
			for _, e := range fresh {
				patches = append(patches, jsonPatch{Op: "add", Path: base + "/-", Value: e})
			}
		}
	}

	return patches
}
