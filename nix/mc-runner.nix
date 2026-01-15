{
  pkgs,
  buildGoModule,
  frontend,
  cloudflaredFlags,
  rev,
  ...
}:
let
  cloudflared-wrapper = ((import ./cloudflared-wrapper.nix) { inherit pkgs; }) cloudflaredFlags;
in
buildGoModule (finalAttrs: {
  src = ./..;
  pname = "mc-runner";
  version = "1.0.0-${rev}";
  vendorHash = "sha256-y7Ou8MvRERtYLt5kqPXD+gbzMMvTF7RuHHKyp4a/hZ8=";

  subPackages = [ "." ];

  disallowedRequisites = [
    frontend
    cloudflared-wrapper
  ];

  doCheck = false;

  preBuild = # sh
    ''
      rm -rf ./frontend
      mkdir -p ./frontend/dist/cloudflared-wrapper

      cp -rv ${frontend}/* ./frontend/dist
      cp -rv ${cloudflared-wrapper}/bin/* ./frontend/dist/cloudflared-wrapper
    '';

  ldflags = [
    "-w"
    "-s"
    "-X github.com/billy4479/mc-runner/internal/config.Version=${finalAttrs.version}"
  ];
})
