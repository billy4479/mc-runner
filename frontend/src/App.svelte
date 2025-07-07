<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { ping, register } from "$lib/api";

  ping();

  function onsubmit(event: Event) {
    event.preventDefault();

    const data = new FormData(event.target);
    register(data.get("token"), data.get("name"));
  }
</script>

<main>
  <h1 class="text-red-300">HEY</h1>
  <Button>A button!</Button>

  <br />

  <form action="" {onsubmit}>
    <input type="text" placeholder="Name" name="name" />
    <input type="text" placeholder="Token" name="token" />
    <input type="submit" value="Register" />
  </form>
</main>
