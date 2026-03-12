{
  self,
  lib,
  config,
  pkgs,
  ...
}:
with lib;
with types;

let
  cfg = config.services.vl-move;
in
{
  options.services = {
    vl-move = {
      enable = mkEnableOption "Enable VL move";
      basePath = mkOption {
        type = path;
      };
      credPath = mkOption {
        type = path;
      };
      timerConfig = mkOption {
        description = ''
          Systemd.timer for this service
          Reference: https://www.freedesktop.org/software/systemd/man/latest/systemd.timer.html
        '';
        type = submodule ({
          freeformType = anything;
          options = {
            OnCalendar = mkOption {
              description = ''
                Systemd.time for this timer
                Reference: https://www.freedesktop.org/software/systemd/man/latest/systemd.time.html
              '';
              type = str;
              example = "*:00:00 Europe/Oslo"; 
            };
          };
        });
      };
    };
  };

  config = {
    systemd.services.vl-move = mkIf cfg.enable ({
      description = "Move files from VL to machine";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      wants = [ "network.target" ];
      script = ''
        ${pkgs.vl.vl-move}/bin/move --local ${cfg.basePath} --credentials ${cfg.credPath}
      '';
    });

    systemd.timers.vl-move = mkIf cfg.enable ({
      description = "Timer for move";
      partOf = [ "vl-move.service" ];
      wantedBy = [ "timers.target" ];
      timerConfig = cfg.timerConfig;
    });
  };
}
