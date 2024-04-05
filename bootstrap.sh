#!/usr/bin/env bash

GOPATH="${HOME}/go"
GOBIN="${GOPATH}/bin"
mkdir -p "${GOBIN}"

ssh_generate_keys() {
  echo "Generating ssh keys..."
  read -rp "Please enter your email: " email
  ssh-keygen -t ed25519 -C "${email}"
}

ssh_add_identity() {
  ssh-add
}

github_upload_public_key() {
  xdg-open https://github.com/settings/keys &>/dev/null
  printf "\n---\n"
  read -rp "Please upload your ssh public key to github.com

Opening https://github.com/settings/keys

Public key:    $(cat "${HOME}/.ssh/id_ed25519.pub")

PRESS ENTER TO CONTINUE"
}


clone_data_repo() {
  read -rp "Please enter your github user: (e.g. 'alexandremahdhaoui') " github_user
  read -rp "Please enter the name of your data repository: (e.g. 'data') " data_repo
  read -rp "Please enter the branch you want to switch to: (e.g. 'main') " data_branch

  DEST_DIR="${GOPATH}/src/github.com/${github_user}/${data_repo}"
  REPO_ADDR="git@github.com:${github_user}/${data_repo}.git"
  mkdir -p "$(dirname "${DEST_DIR}")"

  (
    cd "$(dirname "${DEST_DIR}")" || { echo "Failed changing directory to \"${DEST_DIR}\"" && exit 1 ; }
    git clone "${REPO_ADDR}" || { echo "Failed cloning repository \"${REPO_ADDR}\"" && exit 1 ; }
    cd ./data || exit 1
    git switch "${data_branch}" || { echo "Failed switching to branch \"${data_branch}\"" && exit 1 ; }
  )
}

vib_install() {
  go install github.com/alexandremahdhaoui/vib/cmd/vib@latest
  vib render thiswillfail &>/dev/null # this install vib config.
}

ssh_generate_keys
ssh_add_identity
github_upload_public_key
clone_data_repo
vib_install
