{ pkgs }:
{
  vl-move = pkgs.buildGoModule ({
    version = "0.1.0";
    src = ./.;
    vendorHash = "sha256-o4cvvPREldytDd+7JSH+ewywvnnxpo9zJ8emkyu2i2M=";
    pname = "vl-move";
    subPackages = [ "move" ];
    doCheck = false;
  });

  vl-upload = pkgs.buildGoModule ({
    version = "0.1.0";
    src = ./.;
    vendorHash = "sha256-o4cvvPREldytDd+7JSH+ewywvnnxpo9zJ8emkyu2i2M=";
    pname = "vl-upload";
    subPackages = [ "upload" ];
    doCheck = false;
  });
}
