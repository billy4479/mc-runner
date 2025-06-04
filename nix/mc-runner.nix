{ buildGoModule, frontend, ... }:
buildGoModule {
  src = ./..;
  pname = "mc-runner";
  version = "1.0.0";
  vendorHash = "sha256-IHQgR0IXsaKX+2vOdjDuw7lXQfg8kFaIRpzbWk+1whA=";

  buildInputs = [ frontend ];

  ldflags = [
    "-w"
    "-s"
    "-X github.com/billy4479/mc-runner/internal.FrontendPath=${frontend}"
    "-X github.com/billy4479/mc-runner/internal.BuildMode=release"
  ];
}
