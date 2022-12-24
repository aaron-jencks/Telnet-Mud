CC=go
GOBJECTS=player_service.o inventory_service.o room_service.o logger.o csv.o db.o entities.o player_controller.o server_utils.o room_controller.o crud.o item_service.o transition_service.o game_controller.o command_service.o command_controller.o transition_controller.o item_controller.o inventory_controller.o
TEST_OBJECTS=test_defs.o assert.o csv_tests.o db_tests.o

telnet: telnetMain.o telnet.o
	$(CC) -o $@ telnetMain.o telnet.o logger.o telnet_options.o

telnetMain.o: ./tests/telnet.go telnet.o
	$(CC) -o $@ -c $<

telnet.o: ./net/telnet.go logger.o telnet_options.o
	$(CC) -o $@ -c $<

telnet_options.o: ./net/options.go
	$(CC) -o $@ -c $<

all: mud .

mud: main.o $(GOBJECTS)
	$(CC) main.o $(GOBJECTS) -o $@

main.o: main.go room_service.o entities.o player_controller.o server_utils.o room_controller.o game_controller.o command_controller.o transition_controller.o item_controller.o inventory_controller.o
	$(CC) -o $@ -c $<

server_utils.o: utils/server.go logger.o
	$(CC) -o $@ -c $<

%_controller.o: controllers/%.go %_service.o server_utils.o
	$(CC) -o $@ -c $<

game_controller.o: controllers/game.go server_utils.o
	$(CC) -o $@ -c $<

crud.o: utils/crud.go db.o
	$(CC) -o $@ -c $<

%_service.o: services/%.go crud.o entities.o db.o
	$(CC) -o $@ -c $<

transition_service.o: services/transition.go crud.o entities.o db.o command_service.o
	$(CC) -o $@ -c $<

logger.o: ./ui/logger.go
	$(CC) -o $@ -c $<

csv.o: ./io/csv.go logger.o
	$(CC) -o $@ -c $<

db.o: ./io/db.go logger.o csv.o
	$(CC) -o $@ -c $<

entities.o: ./entities.go logger.o db.o
	$(CC) -o $@ -c $<

csv_tests.o: ./tests/csv_tests.go logger.o test_defs.o csv.o assert.o
	$(CC) -o $@ -c $<

db_tests.o: ./tests/db_tests.go logger.o test_defs.o db.o assert.o
	$(CC) -o $@ -c $<

.PHONY: cleanData
cleanData:
	rm -rf ./data

.PHONY: clean
clean:
	rm -rf main.o rcreator.o $(GOBJECTS) ./mud $(TEST_OBJECTS) testRunner.o ./tester

.PHONY: install
install: mud
	cp -R ./data /usr/local/bin/data
	cp ./mud /usr/local/bin/mudServer

.PHONY: uninstall
uninstall:
	rm /usr/local/bin/mudServer

.PHONY: backup
backup:
	mount /dev/fd0 /media/floppy0
	rm -rf /media/floppy0/mudServer
	cp -R ../mudServer /media/floppy0
	umount /dev/fd0

