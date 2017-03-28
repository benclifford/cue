#!/bin/bash -ex

function asdocker {
  sudo $@
}

DIR=$(dirname $0)

pushd $DIR/../dockerfiles
ls

cd $1
asdocker docker build . | tee /tmp/cuebuild.123
# TODO needs to be generated for multithreadedness

popd

IMAGE=$(cat /tmp/cuebuild.123 | tail -n1 | cut --fields=3 --delimiter=" ")

# the above stuff ensures a docker image exists. there are different ways
# in which that might happen - not only purely dockerfile based, but using
# more interesting docker image building stuff.

# we'll run the right image here, but we won't have the right
# permissions: want home directory mounted, and want the current user
# to exist. (and other environmental attributes - X, ssh agent,
# for example)

# TODO: split run into create/prep/run core command/delete

CIDFILE=/tmp/cuebuild.123-cid # TODO needs to be unique
TMPDIR=$(pwd)
ROOTPREPFILE=${TMPDIR}/cuebuild.123-root # TODO needs to be unique
USERPREPFILE=${TMPDIR}/cuebuild.123-user # TODO needs to be unique

export DOCKEROPTS=""

# prepare homedir mount
export DOCKEROPTS="$DOCKEROPTS -v $HOME:$HOME"

asdocker rm -f $CIDFILE
asdocker rm -f $ROOTPREPFILE
asdocker rm -f $USERPREPFILE

# TODO: this rootprepfile should be copied into the container
# after create, so that it can be deleted rather than having
# to exist as long as the container does. (I needed that before,
# but why? maybe for long lived containers that were to be
# restarted?)

echo "#!/bin/bash" > $ROOTPREPFILE
echo "echo Executing root prep file" >> $ROOTPREPFILE

echo "#!/bin/bash" > $USERPREPFILE
echo "echo Executing user prep file" >> $USERPREPFILE



# Create user
# This needs to copy enough of the user database across
# as needed. If we were in LDAP land, we'd point at the
# LDAP server here rather than using useradd.
echo "echo Creating user" >> $ROOTPREPFILE
echo "useradd benc --uid=1000" >> $ROOTPREPFILE

# Give user sudo rights
# (incidentally mess up other sudo defaults)
echo "echo '%sudo   ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers" >> $ROOTPREPFILE

echo "adduser root sudo" >> $ROOTPREPFILE
echo "adduser benc sudo" >> $ROOTPREPFILE

echo "cd $(pwd)" >> $USERPREPFILE # TODO: check this will be mounted and error with useful message if not

echo "/bin/bash" >> $USERPREPFILE

# Run final command
echo "echo Running user prep file" >> $ROOTPREPFILE
echo "sudo -u benc -i $USERPREPFILE" >> $ROOTPREPFILE


chmod a+x $ROOTPREPFILE
chmod a+x $USERPREPFILE

asdocker docker create $DOCKEROPTS --rm -t -i --cidfile=$CIDFILE $IMAGE $ROOTPREPFILE

CID=$(cat $CIDFILE)

asdocker docker start --attach --interactive $CID

