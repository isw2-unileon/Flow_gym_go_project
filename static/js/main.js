const form = document.getElementById("recommendation-form");
const resultDiv = document.getElementById("result");
const machineSlots = document.querySelectorAll(".machine-slot");
const exerciseOptions = document.getElementById("exercise-options");
const availableCount = document.getElementById("available-count");
const occupiedCount = document.getElementById("occupied-count");
const availableList = document.getElementById("available-list");

form.addEventListener("submit", async function (event) {
    event.preventDefault();

    const exercise = document.getElementById("exercise").value.trim();
    resultDiv.classList.remove("result-success", "result-error", "result-neutral");

    if (!exercise) {
        resultDiv.classList.add("result-error");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>Please enter an exercise name.</p>
        `;
        return;
    }

    try {
        const response = await fetch(`/recommendation?exercise=${encodeURIComponent(exercise)}`);
        const data = await response.json();

        if (!response.ok) {
            resultDiv.classList.add("result-error");
            resultDiv.innerHTML = `
                <h2>Recommendation</h2>
                <p>${data.message || "Could not get recommendation."}</p>
            `;
            return;
        }

        resultDiv.classList.add("result-success");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p><strong>Requested Exercise:</strong> ${data.requested_exercise}</p>
            <p><strong>Recommended Exercise:</strong> ${data.recommended_exercise}</p>
            <p><strong>Muscle Group:</strong> ${data.muscle_group}</p>
            <p><strong>Machine:</strong> ${data.machine}</p>
        `;
    } catch (error) {
        resultDiv.classList.add("result-error");
        resultDiv.innerHTML = `
            <h2>Recommendation</h2>
            <p>An unexpected error occurred while fetching the recommendation.</p>
        `;
    }
});

async function loadMachines() {
    try {
        const response = await fetch("/machines");
        const machines = await response.json();

        const availableMachines = machines.filter(machine => machine.is_available);
        const occupiedMachines = machines.filter(machine => !machine.is_available);

        availableCount.textContent = availableMachines.length;
        occupiedCount.textContent = occupiedMachines.length;

        if (availableMachines.length > 0) {
            availableList.textContent = availableMachines.map(machine => machine.name).join(", ");
        } else {
            availableList.textContent = "No machines available";
        }

        machineSlots.forEach(slot => {
            slot.classList.remove("available", "occupied");

            const machineName = slot.dataset.machineName;
            const machine = machines.find(m => m.name === machineName);

            const statusText = slot.querySelector(".machine-status");
            if (statusText) statusText.textContent = "Unknown";

            if (machine) {
                slot.dataset.machineId = machine.id;
                slot.dataset.available = machine.is_available;

                if (machine.is_available) {
                    slot.classList.add("available");
                    if (statusText) statusText.textContent = "Available";
                } else {
                    slot.classList.add("occupied");
                    if (statusText) statusText.textContent = "Occupied";
                }
            }
        });
    } catch (error) {
        console.error("Could not load machines:", error);
    }
}

async function toggleMachineAvailability(machineId, currentAvailability) {
    try {
        const newAvailability = !currentAvailability;

        const response = await fetch(
            `/machines/update-availability?id=${machineId}&available=${newAvailability}`
        );

        if (!response.ok) {
            console.error("Could not update machine availability");
            return;
        }

        await loadMachines();
    } catch (error) {
        console.error("Error updating machine availability:", error);
    }
}

machineSlots.forEach(slot => {
    slot.addEventListener("click", async function () {
        const machineId = this.dataset.machineId;
        const currentAvailability = this.dataset.available === "true";

        if (!machineId) {
            return;
        }

        await toggleMachineAvailability(machineId, currentAvailability);
    });
});

async function loadExercises() {
    try {
        const response = await fetch("/exercises");
        const exercises = await response.json();

        exerciseOptions.innerHTML = "";

        exercises.forEach(exercise => {
            const option = document.createElement("option");
            option.value = exercise.name;
            exerciseOptions.appendChild(option);
        });
    } catch (error) {
        console.error("Could not load exercises:", error);
    }
}

loadMachines();
loadExercises();