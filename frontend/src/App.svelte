<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { ping } from "$lib/api";

  ping();
</script>

<main>
  <h1 class="text-red-300">HEY</h1>
  <Button>A button!</Button>
</main>
