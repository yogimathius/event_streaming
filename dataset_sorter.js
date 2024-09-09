const { readFileSync, writeFileSync } = require("fs");

const datasets = [
  "assets/dataset_1.json",
  "assets/dataset_2.json",
  "assets/dataset_3.json",
  "assets/dataset_4.json",
  "assets/dataset_5.json",
];

datasets.forEach((datasetPath) => {
  const dataset = JSON.parse(readFileSync(datasetPath, "utf8"));

  const sortedDataset = dataset.sort((a, b) => {
    const [hoursA, minutesA] = a.timestamp.split(":").map(Number);
    const [hoursB, minutesB] = b.timestamp.split(":").map(Number);

    const totalMinutesA = hoursA * 60 + minutesA;
    const totalMinutesB = hoursB * 60 + minutesB;

    console.log(totalMinutesA, totalMinutesB);
    return totalMinutesA - totalMinutesB;
  });

  const reIdDataset = sortedDataset.map((data, index) => {
    return { ...data, id: index + 1 };
  });

  writeFileSync(
    datasetPath.replace(".json", "_sorted.json"),
    JSON.stringify(reIdDataset, null, 2)
  );
});
