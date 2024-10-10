export default {
  fetch() {
    return new Response("hello");
  },
} satisfies ExportedHandler<Env>;
