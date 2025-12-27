{
  dockerTools,
  mc-runner,
  frontend,
  ...
}:

dockerTools.buildLayeredImage {
  name = "mc-runner";
  tag = "latest";

  contents = [
    dockerTools.caCertificates
    mc-runner
    frontend
  ];

  config = {
    Cmd = [ "/bin/mc-runner" ];
  };
}
