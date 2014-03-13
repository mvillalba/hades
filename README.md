Hades
=====
[![Build Status](https://travis-ci.org/mvillalba/hades.png?branch=dev)](https://travis-ci.org/mvillalba/hades) [![GoDoc](https://godoc.org/github.com/mvillalba/hades?status.png)](https://godoc.org/github.com/mvillalba/hades)

_NOTE: UNDER DEVELOPMENT. NOT READY FOR USE._

Hades is a generic software licensing platform written in Go. It supports
various licensing schemes to suit various needs.

The supported (or planned) licensing schemes are:
 * Product key: a key is generated and validated by the licensed software, no
   server is required. This does not prevent multiple installs or allows the
   licensed software to derive much information from the license key itself.
 * Product key with online activation: same as before but the key is activated
   online and the licensed software can download license restriction and
   ownership data from the server. Activation/validation can be used to prevent
   multiple-installs to a point.
 * Product key with online activation and authorization: as above, but the
   licensed software requires a timed authorization from the license server to
   run for a given period of time or until a specific point in time.
   Authorizations can be conditioned to a number of running instances.
 * Same as the above options, but with a license file instead of a bare key:
   the file contains all unmutable data associated with the license
   (restrictions, owner, expiration, etc.) and possibly SSL certificates (for
   interactions with the licensing server), and possibly miscellenous
   license-related files.

Secondary license schemes (can be combined with any of the above):
 * Quota-based license: A certain quota of a given resource (image files
   processed, Web pages served, MB stored, etc.) is associated with the license
   and directly monitored by the remote license server. This can be used to
   have multiple instances of the software running but limiting the total
   number of operations performed or similar. Quotas can be configured to reset
   after a certain period of time has elapsed, be based on a “live” counter
   that goes down when the instance making use of those resources releases them
   or dies, or be fixed (the license is issued for a one-off quota of X
   resources).
 * Machine-lock: Lock a license to a given machine by MAC address or network
   hostname. Requires an online-enabled license.

Hades consists of the following components:
 * Server: manages all licenses for one or more products, handles
   authorizations, activations, and validations as well as reporting, and
   license expiration/renewals.
 * Server CLI: allows administration of the server (pushing, revoking,
   renewing, fetching, and updating licenses).
 * Client libraries: these are linked against licensed software and the handle
   license validation and communication with the licensing server, if required.
 * License utility: generates license keys and files without a server using the
   same libraries the server would use. Licenses generated in this manner can
   then be sent to the license server, if so desired.
