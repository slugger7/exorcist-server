#!/usr/bin/env bash
echo "Generating mocks"

REPO_DIR="./internal/repository"
SERVICE_DIR="./internal/service"

MOCK_DIR="./internal/mock"
MOCK_REPO_DIR="${MOCK_DIR}/repository"
MOCK_SERVICE_DIR="${MOCK_DIR}/service"

mkdir -p ${MOCK_DIR}

echo "Generating repository mocks"
mkdir -p ${MOCK_REPO_DIR}
mockgen -source=${REPO_DIR}/repository.go > ${MOCK_REPO_DIR}/repository.go

mkdir -p ${MOCK_REPO_DIR}/job
mockgen -source=${REPO_DIR}/job/job.go >  ${MOCK_REPO_DIR}/job/job.go

mkdir -p ${MOCK_REPO_DIR}/library
mockgen -source=${REPO_DIR}/library/library.go >  ${MOCK_REPO_DIR}/library/library.go

mkdir -p ${MOCK_REPO_DIR}/library_path
mockgen -source=${REPO_DIR}/library_path/library_path.go >  ${MOCK_REPO_DIR}/library_path/library_path.go

mkdir -p ${MOCK_REPO_DIR}/user
mockgen -source=${REPO_DIR}/user/user.go >  ${MOCK_REPO_DIR}/user/user.go

mkdir -p ${MOCK_REPO_DIR}/video
mockgen -source=${REPO_DIR}/video/video.go >  ${MOCK_REPO_DIR}/video/video.go

mkdir -p ${MOCK_REPO_DIR}/image
mockgen -source=${REPO_DIR}/image/image.go > ${MOCK_REPO_DIR}/image/image.go

echo "Generate service mocks"
mkdir -p ${MOCK_SERVICE_DIR}
mockgen -source=${SERVICE_DIR}/service.go > ${MOCK_SERVICE_DIR}/service.go

mkdir -p ${MOCK_SERVICE_DIR}/library
mockgen -source=${SERVICE_DIR}/library/library.go > ${MOCK_SERVICE_DIR}/library/library.go

mkdir -p ${MOCK_SERVICE_DIR}/library_path
mockgen -source=${SERVICE_DIR}/library_path/library_path.go > ${MOCK_SERVICE_DIR}/library_path/library_path.go

mkdir -p ${MOCK_SERVICE_DIR}/user
mockgen -source=${SERVICE_DIR}/user/user.go > ${MOCK_SERVICE_DIR}/user/user.go

mkdir -p ${MOCK_SERVICE_DIR}/video
mockgen -source=${SERVICE_DIR}/video/video.go > ${MOCK_SERVICE_DIR}/video/video.go

mkdir -p ${MOCK_SERVICE_DIR}/image
mockgen -source=${SERVICE_DIR}/image/image.go > ${MOCK_SERVICE_DIR}/image/image.go

echo "Mocks generated"
