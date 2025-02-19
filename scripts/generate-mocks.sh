#!/usr/bin/env bash
echo "Generating mocks"

MOCK_DIR="./internal/mock"
REPO_DIR="./internal/repository"

mkdir -p ${MOCK_DIR}

MOCK_REPO_DIR="${MOCK_DIR}/repository"
mkdir -p ${MOCK_REPO_DIR}
mockgen -source=${REPO_DIR}/repository.go > ${MOCK_DIR}/repository/repository.go

mkdir -p ${MOCK_REPO_DIR}/job
mockgen -source=${REPO_DIR}/job/job.go > ${MOCK_DIR}/repository/job/job.go

mkdir -p ${MOCK_REPO_DIR}/library
mockgen -source=${REPO_DIR}/library/library.go > ${MOCK_DIR}/repository/library/library.go

mkdir -p ${MOCK_REPO_DIR}/library_path
mockgen -source=${REPO_DIR}/library_path/library_path.go > ${MOCK_DIR}/repository/library_path/library_path.go

mkdir -p ${MOCK_REPO_DIR}/user
mockgen -source=${REPO_DIR}/user/user.go > ${MOCK_DIR}/repository/user/user.go

mkdir -p ${MOCK_REPO_DIR}/video
mockgen -source=${REPO_DIR}/video/video.go > ${MOCK_DIR}/repository/video/video.go
