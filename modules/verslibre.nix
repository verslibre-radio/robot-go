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
  cfg = config.services.verslibre;
in
{
  options.services = {
    verslibre = {
      enable = mkEnableOption "Enable the VL robot";
      basePath = mkOption {
        type = path;
        default = "/var/lib/robot";
      };
      credPath = mkOption {
        type = path;
        default = "/var/lib/robot/cred.json";
      };
      archivePath = mkOption {
        type = path;
        default = "/var/lib/robot/archive";
      };
      dbPath = mkOption {
        type = path;
        default = "/var/lib/robot/metadata.sql";
      };
      timerConfig = mkOption {
        default = { };
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
              default = "Hourly";
            };
          };
        });
      };
    };
  };

  config = {
    systemd.services.verslibre = mkIf cfg.enable ({
      description = "Move and upload files";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      wants = [ "network.target" ];
      serviceConfig.EnvironmentFile="/etc/vl-upload.env";
      script = ''
        echo "Starting moving of files"
        ${pkgs.vl.vl-move}/bin/move --local ${cfg.basePath} --credentials ${cfg.credPath}

        echo "Delete files older than 14 days"
        find ${cfg.archivePath} -type f -mtime +7 -delete

        echo "Starting upload"
        ${pkgs.vl.vl-upload}/bin/upload --local ${cfg.basePath} --credentials ${cfg.credPath} --metadata ${cfg.dbPath} --archive ${cfg.archivePath}

        echo "Finished the program"
      '';
    });

    systemd.timers.verslibre = mkIf cfg.enable ({
      description = "Timer for Verslibre";
      partOf = [ "verslibre.service" ];
      wantedBy = [ "timers.target" ];
      timerConfig = cfg.timerConfig;
    });
  };
}
