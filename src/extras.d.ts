// This is how workerd imports WebAssembly
declare module "*.wasm" {
  const self: WebAssembly.Module;
  export default self;
}

// We need an ImageData definition to satisfy jSquash
type ImageData = unknown;
