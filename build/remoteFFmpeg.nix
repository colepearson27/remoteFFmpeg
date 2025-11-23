{
  buildGoModule,
  fetchFromGitHub,
}:

buildGoModule rec {
  pname = "remoteFFmpeg";
  version = "9f9c350";

  # Use the vender hash below to get the sha-256 used in the fetchFromGithub
  # vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";

  vendorHash = null;

  src = fetchFromGitHub {
    owner = "Robotboy26";
    repo = "remoteFFmpeg";
    rev = version;
    sha256 = "sha256-tLIj4v5n1zeiSGIXQit5UIU4q7l4KKFprZyAbLs3IRk=";
  };
}
