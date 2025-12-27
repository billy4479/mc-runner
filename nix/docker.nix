{
  dockerTools,
  mc-runner,
  mc-java,
  ...
}:

dockerTools.buildLayeredImage {
  name = "mc-runner";
  tag = "latest";

  contents = [
    dockerTools.caCertificates
    mc-runner
    mc-java
  ];

  config = {
    Cmd = [ "/bin/mc-runner" ];
  };
}
