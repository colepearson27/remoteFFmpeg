{
  buildGoModule,
  fetchFromGitHub,
}:

buildGoModule rec {
  pname = "remoteFFmpeg";
  version = "dfc1172";

  # Use the vender hash below to get the sha-256 used in the fetchFromGithub
  # vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";

  vendorHash = null;

  src = fetchFromGitHub {
    owner = "Robotboy26";
    repo = "remoteFFmpeg";
    rev = version;
    sha256 = "sha256-QuRvzV8sc6IRvzCcvdU7gLwUe4ywgRmG4XNpcU4nxUI=";
  };
}
