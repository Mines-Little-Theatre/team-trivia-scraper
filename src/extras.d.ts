// This is how workerd imports WebAssembly
declare module "*.wasm" {
  const self: WebAssembly.Module;
  export default self;
}

// ImageData isn't included in ECMAScript, but workerd has it
type ImageData = unknown;
