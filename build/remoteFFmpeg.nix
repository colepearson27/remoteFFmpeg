{
  buildGoModule,
  fetchFromGitHub,
}:

buildGoModule rec {
  pname = "remoteFFmpeg";
  version = "41ebe68";

  # Use the vender hash below to get the sha-256 used in the fetchFromGithub
  # vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";

  vendorHash = null;

  src = fetchFromGitHub {
    owner = "Robotboy26";
    repo = "remoteFFmpeg";
    rev = version;
    sha256 = "sha256-Zym25yrBO3A9PT5cKrI5L1VAgRGC6ZRKtyE46TAIbwE=";
  };
}
