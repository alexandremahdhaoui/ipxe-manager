package main

// TODO: implement me:
//  - The point of this binary is to set up a local container registry onto which we may push container images.
//  - Once images are pushed to the registry they can be used in the testenv and for chart-testing.
//  - Finally, the binary should also take care of cleaning up the registry.
//
// Consideration: should the local-container-registry run as a container in the default namespace we must ensure
// connectivity between kind pods and the registry.
