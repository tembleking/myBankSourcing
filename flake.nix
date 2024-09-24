{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    let
      overlays.default = final: prev: {
        clerk = prev.callPackage ./clerk.nix { };
        clerkd = prev.callPackage ./clerkd.nix { };
      };
      flake = flake-utils.lib.eachDefaultSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
            overlays = [ self.overlays.default ];
          };
        in
        {
          packages = {
            default = pkgs.clerk;
            inherit (pkgs) clerk clerkd;
          };
          devShells.default =
            with pkgs;
            mkShell {
              packages = [
                # Add here dependencies for the project.
                just
                ginkgo
                go
                mockgen
                buf
                golangci-lint
                gofumpt
                go-migrate
                gorm-gentool
                govulncheck
              ];

              inputsFrom = [
                clerk
                clerkd
              ];
            };

          formatter = pkgs.nixfmt-rfc-style;
        }
      );
    in
    flake // { inherit overlays; };
}
