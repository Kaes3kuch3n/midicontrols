<script lang="ts">
    import {
        GetActions,
        GetMIDIDevices,
        ListenForInput,
        LoadSettings,
        SelectDevice,
        SetCommand
    } from "../wailsjs/go/gui/App";
    import { onMount } from "svelte";

    let devices = [];
    let selectedDevice = undefined;

    let actions;

    let inputEvent;
    let release = false;
    let command = "";

    $: recordedInput = inputEvent ? `${eventTypeToString(inputEvent.type)}, Value: ${inputEvent.value}` : 'None';


    onMount(async () => {
        const settings = await LoadSettings();
        actions = settings.actions;
        devices = await GetMIDIDevices();
        await selectDevice(settings.selectedDevice);
    });

    async function listen() {
        inputEvent = await ListenForInput();
        console.log(inputEvent);
    }

    function eventTypeToString(eventType) {
        switch (eventType) {
            case 0:
                return "Note";
            case 1:
                return "CC";
            case 2:
                return "Prog";
        }
    }

    async function selectDevice(device) {
        if (!device && !selectedDevice) return;
        selectedDevice = await SelectDevice(device ?? selectedDevice);
        console.log(selectedDevice);
    }

    async function setCommand() {
        await SetCommand(inputEvent, release, command);
        actions = await GetActions();
    }
</script>

<main>
    <select name="device" id="device" bind:value={selectedDevice} on:change={selectDevice}>
        {#each devices as device}
            <option value="{device}">{device}</option>
        {/each}
    </select>
    <section id="set-command">
        <button on:click={listen}>Record</button>
        <p>{recordedInput}</p>
        <label>
            Command
            <input type="text" name="command" id="command" bind:value={command}>
        </label>
        <label>
            <input type="checkbox" name="release" id="release" bind:value={release}>
            Release?
        </label>
        <button on:click={setCommand}>Save</button>
    </section>
    {#if (actions)}
        <table>
            <thead>
            <tr>
                <th>Type</th>
                <th>Value</th>
                <th>Press</th>
                <th>Release</th>
                <th></th>
            </tr>
            </thead>
            <tbody>
            {#each Object.entries(actions["0"]) as [key, action]}
                <tr>
                    <td>Note</td>
                    <td>{key}</td>
                    <td>{action.press ?? ""}</td>
                    <td>{action.release ?? ""}</td>
                    <td>x</td>
                </tr>
            {/each}
            {#each Object.entries(actions["1"]) as [key, action]}
                <tr>
                    <td>CC</td>
                    <td>{key}</td>
                    <td>{action.press ?? ""}</td>
                    <td>{action.release ?? ""}</td>
                    <td>x</td>
                </tr>
            {/each}
            {#each Object.entries(actions["2"]) as [key, action]}
                <tr>
                    <td>Prog</td>
                    <td>{key}</td>
                    <td>{action.press ?? ""}</td>
                    <td>{action.release ?? ""}</td>
                    <td>x</td>
                </tr>
            {/each}
            </tbody>
        </table>
    {/if}
</main>

<style></style>
