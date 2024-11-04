#!/bin/bash
AMOUNT_TESTS=100

echo "1" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGet$  roarindex && \
echo "2" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultiple$  roarindex && \
echo "3" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultipleDifferentMaps$  roarindex && \
echo "4" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultipleSingleMapGrow$  roarindex && \
echo "5" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultipleMultiMapGrow$  roarindex && \
echo "6" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultipleMultiMapGrowSameValues$  roarindex && \
echo "7" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexSetAndGetMultipleMultiMapGrowDiffValues$  roarindex && \
echo "8" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexLargeSetAndGet$  roarindex && \
echo "9" && \
go test -timeout 30m -count $AMOUNT_TESTS  -run ^TestRoarIndexLargeSetAndGetAutoGrow$  roarindex  