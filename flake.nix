{
  description = "Move VL";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
    }@inputs:
    let
      inherit (builtins) attrValues;

      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        overlays = [
          (final: prev: {
            vl = import ./packages/default.nix { pkgs = final; };
          })
        ];
      };

      treefmt = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;

    in
    {
      overlays.vl = (
        final: prev: {
          vl = final.callPackage ./packages/default.nix { };
        }
      );

      formatter.${system} = treefmt.config.build.wrapper;

      nixosModules = {
        verslibre = ./modules/verslibre.nix;
      };

      packages.${system} = {
        inherit (pkgs.vl)
          vl-move
          vl-upload
          ;
      };

      checks.${system} = {
        formatting = treefmt.config.build.check self;
      };

      devShells.${system} = {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
          shellHook = ''
            source .env
            root=$(git rev-parse --show-toplevel)
          '';
        };
      };
    };
}
