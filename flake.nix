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
      rec {
        packages = rec {
          frontend = callPackage (import ./nix/frontend.nix) { };
          mc-runner = callPackage (import ./nix/mc-runner.nix) {
            inherit frontend;
            cloudflaredFlags = "version";
            rev = self.shortRev or self.dirtyShortRev or "dirty";
          };
          mc-java = callPackage (import ./nix/java.nix) { };

          cloudflared-wrapper = ((import ./nix/cloudflared-wrapper.nix) { inherit pkgs; }) "version";

          default = mc-runner;
        };

        devShells.default = (import ./nix/shell.nix (pkgs // packages));
      }
    );
}
