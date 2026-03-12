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
  cfg = config.services.vl-upload;
in
{
  options.services = {
    vl-upload = {
      enable = mkEnableOption "Enable VL upload";
      basePath = mkOption {
        type = path;
        default = "/var/lib/robot";
      };
      archivePath = mkOption {
        type = path;
        default = "/var/lib/robot/archive";
      };
      credPath = mkOption {
        type = path;
        default = "/var/lib/robot/cred.json";
      };
      dbPath = mkOption {
        type = path;
        default = "/var/lib/robot/metadata.sql";
      };
      timerConfig = mkOption {
        default = {};
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
              default = "*:30:00 Europe/Oslo"; 
            };
          };
        });
      };
    };
  };

  config = {
    systemd.services.vl-upload = mkIf cfg.enable ({
      description = "Upload files to the various places";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      wants = [ "network.target" ];
      serviceConfig.EnvironmentFile="/etc/vl-upload.env";
      script = ''
        ${pkgs.vl.vl-upload}/bin/upload --local ${cfg.basePath} --credentials ${cfg.credPath} --metadata ${cfg.dbPath} --archive ${cfg.archivePath}
      '';
    });

    systemd.timers.vl-upload = mkIf cfg.enable ({
      description = "Timer for upload";
      partOf = [ "vl-upload.service" ];
      wantedBy = [ "timers.target" ];
      timerConfig = cfg.timerConfig;
    });
  };
}
