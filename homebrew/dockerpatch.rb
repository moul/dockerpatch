class Dockerpatch < Formula
  desc "Debug shell streams"
  homepage "https://github.com/moul/dockerpatch"
  url "https://github.com/moul/dockerpatch/archive/v1.0.0.tar.gz"
  sha256 "ca04d88702576c5c3b6e4a97c19f3a605126454e8974b0461b6c296344d4ede6"
  head "https://github.com/moul/dockerpatch.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["GOBIN"] = buildpath
    (buildpath/"src/github.com/moul/dockerpatch").install Dir["*"]
    system "go", "get", "github.com/moul/dockerpatch/cmd/dockerpatch"
    system "go", "build", "-o", "#{bin}/dockerpatch", "-v", "github.com/moul/dockerpatch/cmd/dockerpatch"
  end
  test do
    system "test", "-f", "#{bin}/dockerpatch"
  end
end
