// This is how workerd imports WebAssembly
declare module "*.wasm" {
  const mod: WebAssembly.Module;
  export default mod;
}

// We need an ImageData definition to satisfy jSquash
// It's available in workerd, just not listed in the type definitions :/
type ImageData = unknown;
