package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"net/http"
)

func HandlePodMutate(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	err = ioutil.WriteFile("/tmp/request", body, 0644)
	if err != nil {
		panic(err.Error())
	}

	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("could not deserialize request: %v", err)
	} else if admissionReviewReq.Request == nil {
		w.WriteHeader(http.StatusBadRequest)
		errors.New("malformed admission review: request is nil")
	}

	fmt.Printf("Type: %v \t Event: %v \t Name: %v \n",
		admissionReviewReq.Request.Kind,
		admissionReviewReq.Request.Operation,
		admissionReviewReq.Request.Name,
	)

	var deployment appv1.Deployment

	err = json.Unmarshal(admissionReviewReq.Request.Object.Raw, &deployment)

	if err != nil {
		fmt.Errorf("could not unmarshal deployment on admission request: %v", err)
	}

	var patches []PatchOperation

	// 定义一个init容器

	cmd := []string{"/bin/sh"}
	initContainer := apiv1.Container{
		Name:    "init",
		Image:   "newharbor.ipinyou.com/deepzero/volvo-sso-jar/build:v1.0", // TODO 通过标签获取镜像名称，而不是写死
		Command: cmd,
		Args:    []string{"-c", "mkdir -p /init/agent && cp -r /data/email-volvohuawei-prod/*.jar /init/agent"},
		VolumeMounts: []apiv1.VolumeMount{
			{
				Name:      "init",
				MountPath: "/init/agent",
			},
		},
	}

	// 添加共享目录卷
	patches = append(patches, PatchOperation{
		Op:   "add",
		Path: "/spec/template/spec/volumes",
		Value: []apiv1.Volume{
			{
				Name: "init",
				VolumeSource: apiv1.VolumeSource{
					EmptyDir: &apiv1.EmptyDirVolumeSource{},
				},
			},
		},
	})

	// 添加主容器挂载共享目录
	patches = append(patches, PatchOperation{
		Op:   "add",
		Path: "/spec/template/spec/containers/0/volumeMounts",
		Value: []apiv1.VolumeMount{
			{
				Name:      "init",
				MountPath: "/init/agent",
			},
		},
	})

	// 添加初始化容器
	patches = append(patches, PatchOperation{
		Op:    "add",
		Path:  "/spec/template/spec/initContainers",
		Value: []apiv1.Container{initContainer},
	})

	patchBytes, err := json.Marshal(patches)

	if err != nil {
		fmt.Errorf("could not marshal JSON patch: %v", err)
	}

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
		},
	}

	admissionReviewResponse.Response.Patch = patchBytes

	bytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		fmt.Errorf("marshaling response: %v", err)
	}

	w.Write(bytes)

}
