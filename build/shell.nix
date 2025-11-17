# This is required for impure linking reasons but will not take up much space on the hard drive
# nix-channel --add https://nixos.org/channels/nixos-25.05 nixpkgs
# nix-channel --update

let
  nixpkgs = fetchGit {
      url = "https://github.com/NixOS/nixpkgs/";
      ref = "release-25.05";
  };
  pkgs = import nixpkgs {};
  remoteFFmpeg = pkgs.callPackage ./remoteFFmpeg.nix { };
in

pkgs.mkShellNoCC { # Make a shell without a c compiler (Use mkShell to get a c compiler)
  packages = [ # Package to include in the shell
      pkgs.cowsay
      pkgs.ffmpeg
      remoteFFmpeg
  ];

  shellHook = ''
      echo ${pkgs.cowsay}
      cowsay "Welcome to the remoteFFmpeg build testing env"
      echo ${remoteFFmpeg}
      '';
}
