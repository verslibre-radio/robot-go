{
  description = "Move VL";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs =
    {
      self,
      nixpkgs,
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


    in
    {
      overlays.vl = (final: prev: {
          vl = final.callPackage ./packages/default.nix {};
        });

      nixosModules = {
        verslibre = ./modules/verslibre.nix;
      };

      packages.${system} = {
        inherit (pkgs.vl)
          vl-move
          vl-upload
          ;
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
