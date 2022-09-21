
#include <Arduino.h>
#define ENDSTOP 6
int idx = 0;
int stepDelay = 450; //delay between step signals send to the stepper motor
int maxSteps = 8160; // max steps before the carriage runs into the endstop
int curstepts = 0; // current carriage position

unsigned long mil; // current time
bool measure;

/**
 * step one step into the currently selected direction
 */
void stepF()
{
  if (ENDSTOP == LOW || curstepts > maxSteps) // ensure that the carriage does not run into the endstop or beyond the 0 position
  {
    return;
  }
  
  digitalWrite(3, HIGH);
  delayMicroseconds(stepDelay);
  digitalWrite(3, LOW);
  delayMicroseconds(stepDelay);
}

void setup()
{
  Serial.begin(9600);
  // put your setup code here, to run once:
  pinMode(2, OUTPUT);
  pinMode(3, OUTPUT);
  pinMode(4, OUTPUT);
  pinMode(ENDSTOP, INPUT);
  digitalWrite(2, LOW);
  digitalWrite(4, HIGH);
}

/**
 * find the 0 position of the carriage
 */
void homing()
{
  Serial.println("Homing");
  Serial.println(digitalRead(6));
  digitalWrite(4, HIGH);
  while (digitalRead(6) == HIGH)
  {
    stepF();
  }
  curstepts = 0;
  return;
}

/**
 * move the carriage to the specified position
 */
void goTo(int position)
{
  int amtSteps = position - curstepts;
  int dir = 1;
  if (amtSteps < 0)
  {
    digitalWrite(4, HIGH);
    dir = -1;
  }else{
    digitalWrite(4, LOW);
  }
  amtSteps = abs(amtSteps);
  Serial.println(amtSteps);

  while(amtSteps >0){
    stepF();
    amtSteps --;
    curstepts += dir;
  }


}

void handleInput(){

    String input = Serial.readString();
    input.trim();

    if (input == "c")
    {
      homing();
      delay(50);
      Serial.println("finished homing");
    }
    else if (input.startsWith("to"))
    {
      int position = input.substring(2).toInt();
      goTo(position);
      Serial.print("finished: ");
      Serial.println(position);
    }
    else
    {
      // prints the received integer
      Serial.print("I received: ");
      Serial.println(input);
      idx = input.toInt();
      measure = true;
      mil = millis();
      Serial.print("currentPos: ");
      Serial.println(idx + curstepts);
    }
}

void loop()
{
  if (Serial.available() > 0){
    handleInput();
  }
  



  if (curstepts + idx < maxSteps && curstepts >= 0)
  {
    if (idx > 0)
    {
      digitalWrite(4, LOW);
      stepF();
      idx--;
      curstepts++;
    }
    if (idx == 0 and measure)
    {
      Serial.println("finished");
      measure = false;
    }
    else if (idx < 0)
    {
      digitalWrite(4, HIGH);
      stepF();
      idx++;
      curstepts--;
    }
  }
}