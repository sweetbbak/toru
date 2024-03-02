{
  lib,
  fetchFromGitHub,
  buildGoModule,
}:

buildGoModule rec {
  pname = "toru";
  version = "0.1";

  src = fetchFromGitHub {
    owner = "sweetbbak";
    repo = "toru";
    rev = "c5f6134b42d6dc8fd3f52e598f6fc6b1eae74888";
    hash = "sha256-GiptGi2OZy6eiDa/g8GwPVZpDd8QUxZuzaEy+9J+v9I=";
  };

  # vendorHash = lib.fakeHash;
  /*
  vendorHash = "sha256-alC4/2wTbjJYWGzTDTgQweOicN3xSqfnncok/j16+0E=";
  */
  vendorHash = "sha256-72Vjrl/v4rKojv3BeACcEiOyWK4c9KCyzxJtfNx/7UE=";

  CGO_ENABLED = 0;
  ldflags = ["-s" "-w"];

  tags = ["torrent" "bittorrent" "anime"];
  # proxyVendor = true;
  #
  # buildPhase = ''
  #     go mod vendor
  #     go build -o toru ./cmd/toru
  # '';
  #
  # installPhase = ''
  #     mkdir -p $out/bin
  #     mv toru $out/bin
  # '';

  meta = with lib; {
    homepage = "https://github.com/sweetbbak/toru";
    description = "Bittorrent streaming CLI tool. Stream anime torrents, real-time with no waiting for downloads";
    license = licenses.mit;
    maintainers = with maintainers; [sweetbbak];
    mainProgram = "toru";
  };
}
