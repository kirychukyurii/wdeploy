#!/bin/sh
set -e

WDEPLOY_VERSION="0.0.3"
WANSIBLE_REPOSITORY=$(echo "$HOME/wansible")
DRY_RUN=${DRY_RUN:-}
while [ $# -gt 0 ]; do
	case "$1" in
		--dry-run)
			DRY_RUN=1
			;;
		--*)
			echo "Illegal option $1"
			;;
	esac
	shift $(( $# > 0 ? 1 : 0 ))
done

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

is_dry_run() {
	if [ -z "$DRY_RUN" ]; then
		return 1
	else
		return 0
	fi
}

check_architecture() {
    KERNEL_NAME=$(uname -s)
    HARDWARE_NAME=$(uname -m)

    case "$HARDWARE_NAME" in
      aarch64*)
        HARDWARE_NAME="arm64";;
    esac
}

ansible_install() {
    sh_c='sh -c'
    if is_dry_run; then
    		sh_c="echo"
    fi

    echo "# Executing ansible install script"
    $sh_c "curl -fsSL https://bootstrap.pypa.io/get-pip.py -o $HOME/get-pip.py" > /dev/null 2>&1
    $sh_c "python3 $HOME/get-pip.py --user" > /dev/null 2>&1

	if command_exists ansible; then
		cat >&2 <<-'EOF'
			Warning: the "ansible" command appears to already exist on this system.

			If you already have Ansible installed, this script can cause trouble, which is
			why we're displaying this warning and provide the opportunity to cancel the
			installation.

			If you installed the current Ansible package using this script and are using it
			again to update Ansible, you can safely ignore this message.

			You may press Ctrl+C now to abort this script.
		EOF
		( set -x; sleep 5 )

		$sh_c 'python3 -m pip install --upgrade --user ansible' > /dev/null 2>&1
	else
	    $sh_c 'python3 -m pip install --user ansible' > /dev/null 2>&1
	fi
}

wdeploy_clone_repo() {
    if [ ! -d "$WANSIBLE_REPOSITORY" ]; then
        $sh_c "mkdir $WANSIBLE_REPOSITORY"
    fi

    if [ ! -f "$WANSIBLE_REPOSITORY/playbook.yml" ]; then
        $sh_c "git clone git@github.com:kirychukyurii/wansible.git $WANSIBLE_REPOSITORY" > /dev/null 2>&1
    fi
}

ansible_collections_install() {
    $sh_c "ansible-galaxy collection install -r $WANSIBLE_REPOSITORY/requirements.yml > /dev/null" > /dev/null 2>&1
}

wdeploy_install() {
    $sh_c "curl -fsSL https://github.com/kirychukyurii/wdeploy/releases/download/$WDEPLOY_VERSION/wdeploy_${KERNEL_NAME}_$HARDWARE_NAME.tar.gz -o $HOME/wdeploy_$WDEPLOY_VERSION.tar.gz"
    $sh_c "tar -xf $HOME/wdeploy_$WDEPLOY_VERSION.tar.gz -C $HOME/"
    $sh_c "rm -f $HOME/wdeploy_$WDEPLOY_VERSION.tar.gz"
    $sh_c "chmod +x $HOME/wdeploy"

    echo "wdeploy successfully installed on your system."

	$sh_c "$HOME/wdeploy run --help"
}

# wrapped up in a function so that we have some protection against only getting
# half the file during "curl | sh"
check_architecture
ansible_install
wdeploy_clone_repo
ansible_collections_install
wdeploy_install

