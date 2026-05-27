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
  vendorHash = "sha256-h1QcRSEIFwT7Vhonvjq8sQt7NBT8seTTaKXMDMA0nak=";

  buildInputes = [ frontend ];
  disallowedRequisites = [ frontend ];

  preBuild = # sh
    ''
      rm -rf ./frontend
      mkdir -p frontend
      cp -rv ${frontend} ./frontend/dist
    '';

  ldflags = [
    "-w"
    "-s"
    "-X github.com/billy4479/mc-runner/internal/config.Version=${finalAttrs.version}"
  ];
})
