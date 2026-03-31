class Tuibookie < Formula
  desc "A fast, interactive terminal bookmark manager for CLI commands"
  homepage "https://github.com/orvad/tuibookie"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/orvad/tuibookie/releases/download/v#{version}/tuibookie-darwin-arm64"
      sha256 "PLACEHOLDER_DARWIN_ARM64"

      def install
        bin.install "tuibookie-darwin-arm64" => "tuibookie"
      end
    else
      url "https://github.com/orvad/tuibookie/releases/download/v#{version}/tuibookie-darwin-amd64"
      sha256 "PLACEHOLDER_DARWIN_AMD64"

      def install
        bin.install "tuibookie-darwin-amd64" => "tuibookie"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/orvad/tuibookie/releases/download/v#{version}/tuibookie-linux-arm64"
      sha256 "PLACEHOLDER_LINUX_ARM64"

      def install
        bin.install "tuibookie-linux-arm64" => "tuibookie"
      end
    else
      url "https://github.com/orvad/tuibookie/releases/download/v#{version}/tuibookie-linux-amd64"
      sha256 "PLACEHOLDER_LINUX_AMD64"

      def install
        bin.install "tuibookie-linux-amd64" => "tuibookie"
      end
    end
  end

  test do
    assert_match "tuibookie", shell_output("#{bin}/tuibookie --help 2>&1", 0)
  end
end
