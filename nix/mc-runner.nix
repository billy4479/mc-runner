{ buildGoModule, frontend, ... }:
buildGoModule {
  src = ./..;
  pname = "mc-runner";
  version = "1.0.0";
  vendorHash = "sha256-SUkOb2OGh5xEsUUv5w0YQM7Ctx5DW9R2tBkkz06P2ds=";

  buildInputs = [ frontend ];

  ldflags = [
    "-w"
    "-s"
    "-X github.com/billy4479/mc-runner/internal.FrontendPath=${frontend}"
    "-X github.com/billy4479/mc-runner/internal.BuildMode=release"
  ];
}
