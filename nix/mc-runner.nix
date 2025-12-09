{
  buildGoModule,
  frontend,
  rev,
  ...
}:
buildGoModule (finalAttrs: {
  src = ./..;
  pname = "mc-runner";
  version = "1.0.0-${rev}";
  vendorHash = "sha256-4uAAmzKVbVbe3Ft1QLy87SPUgAgsgevqMj8g4HEPZbs=";

  buildInputs = [ frontend ];

  ldflags = [
    "-w"
    "-s"
    "-X github.com/billy4479/mc-runner/internal.FrontendPath=${frontend}"
    "-X github.com/billy4479/mc-runner/internal.Version=${finalAttrs.version}"
  ];
})
