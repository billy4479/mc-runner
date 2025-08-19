<script lang="ts">
  import { logout } from "$lib/api";
  import { setMeOrError, getMeOrError } from "$lib/state.svelte";
  import { Button } from "flowbite-svelte";
  import LoginScreen from "./LoginScreen.svelte";

  setMeOrError();
</script>

<main class="flex h-screen items-center justify-center gap-4">
  {#if getMeOrError() === null}
    Loading...
  {:else if getMeOrError().error}
    {#if getMeOrError().status == 401}
      <LoginScreen />
    {:else}
      <div>
        <span> An error has occurred </span>
        <pre>{JSON.stringify(getMeOrError(), null, 2)}</pre>
      </div>
    {/if}
  {:else}
    Welcome {getMeOrError().Name.String}

    <Button
      onclick={async () => {
        await logout();
        await setMeOrError();
      }}
    >
      logout
    </Button>
  {/if}
</main>
