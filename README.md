# hpc-cli-go

A tool that provides visibility across an HPC island, cluster network, bare metal cluster, and more. hpc-cli-go can also display broken hosts, and run an optional health check on each host.

## Code geography
```
| File                 | Description                                                                                         |
|----------------------|-----------------------------------------------------------------------------------------------------|
| main.go              | Entry point into the program. Calls hpc.go                                                          |
| hpc.go               | Contains the program logic                                                                          |
| compute.go           | Grouping of functions relating to Compute                                                           |
| baremetalstats.go    | Grouping of functions relating to BareMetalStats (Pulse)                                            |
| computemanagement.go | A group of functions used to interact with the Compute Management API                               |
| filters.go           | Contains all program filter functions                                                               |
| validate.go          | Contains functions to validate input arguments                                                      |
| display.go           | Contains function to write output to terminal and JSON file                                         |
| model.go             | Houses structs and constants used within the tool                                                   |
| go.mod               | Describes the modules properties, including its dependencies on other modules and on versions of Go |
| cli.go               | Struct defining CLI arguments                                                                       |
| version.go           | Maintains the version of the program                                                                |
| tbump.toml           | The current version will be updated when tbump is run (tbump - PyPI)                                |
|____________________________________________________________________________________________________________________________|
```
<br />

## Pre-requisites
```
- AllProxy must be running
- You must be connected to the OCNA VPN
- Your yubikey is required
- Membership of the instance-management group in the Permissions Portal (either directly, or inherited)
```

### Yubikey Modules
A yubikey module such as opensc-pkcs11 or libykcs11 is required. 

By default, the following paths are checked and used if found:

`/Library/OpenSC/lib/opensc-pkcs11.so`\
`/usr/local/lib/opensc-pkcs11.so`\
`/Program Files/OpenSC Project/OpenSC/pkcs11/opensc-pkcs11.dll` \
`/usr/lib/x86_64-linux-gnu/opensc-pkcs11.so` \
`/usr/lib/aarch64-linux-gnu/opensc-pkcs11.so` 

To use another module set the environment variable `OCI_MOLE_MODULE_PATH`
<br />
<br />
### Pin-less Yubikey authentication
By default when the tool is executed, to enable authentication with compute-management, your Yubikey is required. A prompt will be displayed asking for your pin when you execute the tool.
<br />
You can setup a profile in your ~/.oci/config file which will negate the need for you to enter your Yubikey pin. A quick guide on how to do this can be found here https://confluence.oci.oraclecorp.com/pages/viewpage.action?pageId=2747747380
<br/>
<br />
## Getting Started

The package is stored in a private bitbucket repository accessible only via ssh so some configuration is required for Go Modules to work\.

Add bitbucket to the GOPRIVATE environment variable\.

```
export GOPRIVATE=bitbucket.oci.oraclecorp.com
```

Create a URL alias in the git config to use ssh instead of https\.

```
git config --global url."ssh://git@bitbucket.oci.oraclecorp.com:7999/".insteadof "https://bitbucket.oci.oraclecorp.com/scm/"
```
<br />

## Build

```
git clone ssh://git@bitbucket.oci.oraclecorp.com:7999/eng-tool/hpc-cli-go.git

cd hpc-cli-go

go build
```
<br />

## Usage

```
$ ./hpc-cli-go

$ ➜  hpc-cli-go git:(master) ./hpc-cli-go --help
Usage: hpc-cli --region=REGION --version

Flags:
  -h, --help                 Show context-sensitive help.
      --debug                Enable debug mode
  -r, --region=REGION        Example: ad1.us-ashburn-1
  -b, --list-broken-hosts    Returns broken hosts only
  -p, --run-health-check     Runs Pulse health check (ilom/smartnic)
      --version              Show Version

Search Term
  -c, --cluster-network=CLUSTER-NETWORK          Cluster Network OCID e.g ocid1.clusternetwork.oc1.phx.bla
  -d, --hpc-island-id=HPC-ISLAND-ID              HPC Island block ID e.g bldg1-block11
  -m, --bare-metal-cluster=BARE-METAL-CLUSTER    Bare Metal Cluster OCID e.g ocid1.baremetalcluster.oc1.phx.bla
  -s, --rack-id=RACK-ID                          Storekeeper Rack ID e.g sk-blalbla
  -i, --instance-id=INSTANCE-ID                  Instance OCID e.g ocid1.instance.oc1.phx.bla
  -l, --list-all                                 Lists all HPC instances within region
```

<br />

# Examples

1. List all HPC instances within an HPC Island
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --hpc-island-id=bldg3-block3
```
2. List all instances in a Cluster Network
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --cluster-network=ocid1...
```
3. List all instances within a Bare Metal Cluster
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --bare-metal-cluster=ocid1.baremetalcluster.oc1.ap-tokyo-1.anxhiljrwvzy2lycjdax22f5f2cj4gpba5kjyojl723rzqfmwdzwzq5tl4za
```
4. Display HPC instance by Rack ID
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --rack-id=sk-73de362e-a0d8-45f3-bd45-5f63735681cc
```
5. Display HPC instance by Instance OCID
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --instance-id=ocid1.instance.oc1.ap-tokyo-1.anxhiljrwvzy2lyc3dvmlqnhuiwbybpothxk7ysuemq5mwvmhxfysmzug64q
```
6. List all HPC instances within a region
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --list-all
```
7. Add the --list-broken-hosts flag to display broken hosts in conjunction with a search term
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --hpc-island-id=bldg1-block9 --list-broken-hosts
```
8. Run health check on each instance (SMARTNIC and ILOM health)
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --hpc-island-id=bldg3-block3 --run-health-check
```

9. Check health on broken hosts
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --hpc-island-id=bldg3-block3 --list-broken-hosts --run-health-check
```
```
Add the --show-table flag to display the results in a table within the terminal window
```
# Output
## Terminal output
Terminal output is displayed by default when you execute the tool
```
| RACK ID | HPC ISLAND ID | HOST ID | SHAPE | TENANCY OCID | COMPARTMENT OCID | INSTANCE OCID | CLUSTER NETWORK DISPLAY NAME | CLUSTER NETWORK OCID | BMC OCID | INSTANCE POOL | INSTANCE CONFIG | HOST RESERVATION |
```
## JSON
Input
```
./hpc-cli-go --region=ad1.ap-tokyo-1 --instance-id=ocid1.instance.oc1.ap-tokyo-1.anxhiljrwvzy2lyca6qoz3gxte66nw5zalgpcixyblrbbzkf3kqcfhx3ovhq --run-health-check
```
Ouput

  ```json
[
 {
  "rackId": "sk-73de362e-a0d8-45f3-bd45-5f63735681cc",
  "hostId": "sk-a1851a21-cfc2-463a-8b4e-611983912535",
  "instanceId": "ocid1.instance.oc1.ap-tokyo-1.anxhiljrwvzy2lyca6qoz3gxte66nw5zalgpcixyblrbbzkf3kqcfhx3ovhq",
  "compartmentId": "ocid1.compartment.oc1..aaaaaaaaxba2wtxdw763t56o6i2s2gsj4cgfutzfujq4mcpxt3owij7y6rsq",
  "tenancyId": "ocid1.tenancy.oc1..aaaaaaaaqzapziinxpui7ucunahdchlv43qede4ahnq2hhs6sgw7vo7zha2q",
  "hpcIslandId": "bldg3-block3",
  "shape": "x9-2c.36.512",
  "poolName": "compute_standard",
  "bareMetalClusterId": "ocid1.baremetalcluster.oc1.ap-tokyo-1.anxhiljrwvzy2lycwoewy2c7pum3d5hrgu3ezvyf6poufs76ftierbi6cs3q",
  "clusterNetwordId": "ocid1.clusternetwork.oc1.ap-tokyo-1.amaaaaaawvzy2lya63y63ckrizjphlz55dgqezwkm57szjejnyffjswytcha",
  "clusterNetworkDisplayName": "pbscloudiNvp9dT5CysRHNoA-clusternet",
  "instancePoolId": "ocid1.instancepool.oc1.ap-tokyo-1.aaaaaaaaykdjwuzdhx66b2h452y4yjtvghdf6f6sur3hiefxccaemmlr5vhq",
  "instanceConfigurationId": "ocid1.instanceconfiguration.oc1.ap-tokyo-1.aaaaaaaast334zs7xrlhdjqpaepmxi5gsqyuzmm3mcysynyeutrhyxhvphiq",
  "hostReservationId": "ocid1.hostreservation.oc1.ap-tokyo-1.anxhiljrwvzy2lycijalm6qn4nqvhfcma6fsjeoxcrtr2lrv2lfglu766tjq",
  "healthChecks": {
   "hostSmartNicData": {
    "hostDetails": {
     "hostLocation": {
      "building": "nrt3",
      "elevation": "23",
      "rack": "0310",
      "room": "nrt3.1"
     },
     "identity": "2137XCL03T",
     "isPoweredOn": true,
     "platform": "Optimized_Comp_X9-2C_Server.01"
    },
    "smartNicDetailsList": [
     {
      "computeState": "ACTIVE",
      "identity": "NCN2132A08971",
      "ipAddress": "10.2.19.21",
      "isImpaired": false,
      "manufacturer": "Unknown",
      "problemList": [],
      "provisionTime": "2022-11-02T02:57:26Z",
      "smartNicEdacMetrics": {
       "cpuMetricsDetails": {
        "correctableErrorCount": 0,
        "lastObservedTime": "2022-11-07T11:54:16Z",
        "uncorrectableErrorCount": 0
       },
       "memoryMetricsDetails": {
        "correctableErrorCount": 0,
        "lastObservedTime": "2022-11-07T11:54:16Z",
        "uncorrectableErrorCount": 0
       }
      },
      "smartNicInterfaceDetails": [],
      "smartNicLinkFlaps": {
       "count": 0,
       "lastObservedTime": "2022-11-07T11:54:14Z"
      },
      "smartNicReachabilityCheck": {
       "isReachable": true,
       "lastObservedTime": "2022-11-07T11:45:58.531Z"
      },
      "smartNicSelfCheck": {
       "diagnosisDetailsList": [],
       "isSelfCheckPassing": true,
       "lastObservedTime": "2022-11-07T11:54:14Z"
      },
      "uname": "Linux cav-NCN2132A08971 4.14.35-2025.403.1.el7uek.aarch64.emb #2 SMP Fri Nov 13 19:31:08 PST 2020 aarch64 aarch64 aarch64 GNU/Linux",
      "uptime": "2022-11-02T04:27:54Z"
     },
     {
      "computeState": "ACTIVE",
      "identity": "NCN2132A04248",
      "ipAddress": "10.2.19.55",
      "isImpaired": false,
      "manufacturer": "Unknown",
      "problemList": [],
      "provisionTime": "2022-11-02T02:55:50Z",
      "smartNicEdacMetrics": {
       "cpuMetricsDetails": {
        "correctableErrorCount": 0,
        "lastObservedTime": "2022-11-07T11:52:32Z",
        "uncorrectableErrorCount": 0
       },
       "memoryMetricsDetails": {
        "correctableErrorCount": 21,
        "lastObservedTime": "2022-11-07T11:52:32Z",
        "uncorrectableErrorCount": 0
       }
      },
      "smartNicInterfaceDetails": [],
      "smartNicLinkFlaps": {
       "count": 0,
       "lastObservedTime": "2022-11-07T11:52:30Z"
      },
      "smartNicReachabilityCheck": {
       "isReachable": true,
       "lastObservedTime": "2022-11-07T11:55:05.784Z"
      },
      "smartNicSelfCheck": {
       "diagnosisDetailsList": [],
       "isSelfCheckPassing": true,
       "lastObservedTime": "2022-11-07T11:52:30Z"
      },
      "uname": "Linux cav-NCN2132A04248 4.14.35-2025.403.1.el7uek.aarch64.emb #2 SMP Fri Nov 13 19:31:08 PST 2020 aarch64 aarch64 aarch64 GNU/Linux",
      "uptime": "2022-11-02T04:27:53Z"
     }
    ]
   },
   "ilomDetails": {
    "availabilityEvents": [
     {
      "hostSerial": "2137XCL03T",
      "lastPulseReportTime": "2022-11-07T11:51:13.761Z",
      "storeKeeperId": "sk-a1851a21-cfc2-463a-8b4e-611983912535"
     }
    ],
    "hostPowerStatus": {
     "isPoweredOn": true,
     "lastObservedTime": "2022-11-07T11:51:13.761Z"
    },
    "ilomFirmwareVersion": {
     "lastObservedTime": "2022-11-07T11:51:13.761Z",
     "version": "5.0.2.93.a"
    },
    "ilomProcessesDetails": {
     "hasHungProcesses": false,
     "lastObservedTime": "2022-11-07T11:51:13.761Z"
    },
    "ilomReachabilityCheck": {
     "isReachable": true,
     "lastObservedTime": "2022-11-07T11:46:54.684Z"
    },
    "ilomUsers": [
     {
      "isDefaultPassword": false,
      "isLocked": false,
      "lastObservedTime": "2022-11-07T11:51:13.761Z",
      "name": "root"
     },
     {
      "isDefaultPassword": false,
      "isLocked": false,
      "lastObservedTime": "2022-11-07T11:51:13.761Z",
      "name": "ilomusr0"
     },
     {
      "isDefaultPassword": false,
      "isLocked": false,
      "lastObservedTime": "2022-11-07T11:51:13.761Z",
      "name": "ilomusr1"
     },
     {
      "isDefaultPassword": false,
      "isLocked": false,
      "lastObservedTime": "2022-11-07T11:51:13.761Z",
      "name": "hgate"
     }
    ]
   }
  }
 }
]
  ```
