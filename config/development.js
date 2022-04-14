module.exports = {
  indexer: {
    catchup: {
      enabled: true,
      blocksBatchSize: 1,

      blockNumberRanges: [
        [170040, 170040],
        [170498, 170498],
      ],
    },
    realtime: {
      enabled: false,
    },
  },
};
