aws-status
==========

A monitoring tool for AWS components.

### Usage

* Add load balancer and VPN config to a `config.yml` file (see [config.example.yml](config.example.yml))
* Configure AWS credentials, either:
  * On EC2, use an instance profile
  * Export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_REGION`
  * Use `aws configure`
* To run it locally: `make run`
* To build a release: `make`

### License

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
