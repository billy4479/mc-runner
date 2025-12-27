{
  dockerTools,
  mc-runner,
  ...
}:

dockerTools.buildLayeredImage {
  name = "mc-runner";
  tag = "latest";

  contents = [
    dockerTools.caCertificates
    mc-runner
  ];

  config = {
    Cmd = [ "/bin/mc-runner" ];
  };
}
