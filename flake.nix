{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = rec {
          mc-runner = pkgs.callPackage (
            { buildGoModule, ... }:
            buildGoModule {
              src = ./.;
              pname = "mc-runner";
              version = "1.0.0";
              vendorHash = "sha256-P/wg/jogynsre9NBvSkZ6ax8m1LuzTF4Ec84ayin9XY=";
            }
          ) { };
          docker-image = pkgs.dockerTools.buildLayeredImage {
            name = "mc-runner";
            tag = "latest";

            contents = [
              pkgs.dockerTools.caCertificates
              mc-runner
            ];

            config = {
              Cmd = [ "/bin/mc-runner" ];
            };
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            cloudflared
            podman
            turso-cli
            air

            gopls
            golangci-lint
          ];
          nativeBuildInputs = with pkgs; [ go ];
        };
      }
    );
}
