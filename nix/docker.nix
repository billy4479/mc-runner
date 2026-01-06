{
  dockerTools,
  mc-runner,
  mc-java,
  ...
}:

dockerTools.buildImage {
  name = "mc-runner";
  tag = "latest";

  copyToRoot = [
    dockerTools.caCertificates
    mc-runner
    mc-java
  ];

  config = {
    Cmd = [ "/bin/mc-runner" ];
  };
}
