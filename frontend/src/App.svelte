<script lang="ts">
  import { getMe } from "$lib/api";
</script>

<main class="flex h-screen items-center justify-center gap-4">
  {#await getMe()}
    Loading...
  {:then meOrError}
    {#if meOrError.error}
      <div>
        <span> An error has occurred </span>
        <pre>{JSON.stringify(meOrError, null, 2)}</pre>
      </div>
    {:else}
      Welcome {meOrError.Name.String}
    {/if}
  {/await}
</main>
