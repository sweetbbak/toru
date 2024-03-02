{ lib, fetchFromGitHub, buildGoModule }:

buildGoModule rec {
    pname = "toru";
    version = "0.1";

    src = fetchFromGitHub {
      owner = "sweetbbak";
      repo = "toru";
      rev = "747deedd9b19e713f81c208686261e4638e8950e";
      hash = "sha256-iS7M7++r0OejUnqr2st/CXGox1o/DnAQ4w5Kn/nqT7E=";
    };

    # vendorHash = lib.fakeHash;
    vendorHash = "sha256-alC4/2wTbjJYWGzTDTgQweOicN3xSqfnncok/j16+0E=";

    CGO_ENABLED = 1;
    ldflags = [ "-s" "-w" ];

    tags = [ "torrent" "bittorrent" "anime" ];
    proxyVendor = true;

    buildPhase = ''
        go mod vendor
        go build -o toru ./cmd/toru
    '';

    installPhase = ''
        mkdir -p $out/bin
        mv toru $out/bin
    '';

    meta = with lib; {
        homepage = "https://github.com/sweetbbak/toru";
        description = "Stream anime from the command line with the power of torrents";
        license = licenses.mit;
        maintainers = with maintainers; [ sweetbbak ];
        mainProgram = "toru";
    };
}
