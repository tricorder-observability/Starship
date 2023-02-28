TOT	:= $$(git rev-parse --show-toplevel)
DEV_IMAGE_FILE := $(TOT)/tools/docker/dev_image/DEV_IMAGE
PACKER := packer
DEV_IMAGE_REPO := docker.io/tricorderobservability
DEV_IMAGE_TAG := $(shell date "+%Y%m%d%H%M")

.PHONY: build_and_push_ci_image
build_and_push_ci_image:
	$(PACKER) build \
		-var "base_image_and_tag=ubuntu:22.04" \
		-var "generated_image_repository=$(DEV_IMAGE_REPO)/ci_image" \
		-var "generated_image_tag=$(DEV_IMAGE_TAG)" \
		packer_ci_image.json

.PHONY: build_and_push_dev_image
build_and_push_dev_image:
	$(PACKER) build \
		-var "base_image_and_tag=ubuntu:22.04" \
		-var "generated_image_repository=$(DEV_IMAGE_REPO)/dev_image" \
		-var "generated_image_tag=$(DEV_IMAGE_TAG)" \
		packer_dev_image.json
	docker save -o ${DEV_IMAGE_FILE} $(DEV_IMAGE_REPO)/dev_image:$(DEV_IMAGE_TAG)
