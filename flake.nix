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
        inherit (pkgs) callPackage;
      in
      {
        packages = rec {
          frontend = callPackage (import ./nix/frontend.nix) { };
          mc-runner = callPackage (import ./nix/mc-runner.nix) { inherit frontend; };
          docker-image = callPackage (import ./nix/docker.nix) { inherit mc-runner; };
        };

        devShells.default = (import ./nix/shell.nix pkgs);
      }
    );
}
