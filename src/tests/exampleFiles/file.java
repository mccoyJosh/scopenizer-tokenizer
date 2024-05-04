
// Here is a comment :)

/*
 Here
 is
 a
 multiline
 comment
* */

public class file {
    public static void main() {
        int variable = 10;
        System.out.println("hey");
    }

    public static void el(int m) { // 1 (just cause)
        if (m == 0 || m == 2) {   // 1                          ;2
            if (m == 2) {                   // 2        ;4
               System.out.println("m is not");
            } else {                        // 1        ;5
                System.out.println("m is very not");
            }
        } else {                        // 1            ;6
            System.out.println("m is really not");
        }

        if (m > 0 && m % 2 == 0 && boolMethod()) {  // 1 ;7
            System.out.println("m is");
        }
    } // 7 total


    public static void el2(int m) { // 1
        if (m == 0 || m == 2)       // 1        ; 2
            if (m == 2)             // 2            ; 4
                System.out.println("m is not");

        if (m > 0 && m % 2 == 0 && boolMethod())    // 5
            System.out.println("m is");
    }


    public static void el3(int m) { // 1
        if (m == 0 || m >= 2) {      // 1     ;2
            if (m == 2) {           // 2      ;4
                System.out.println("m is not");
            } else if (m == 0) {                //1     ; 5
                System.out.println("m is very not");
            } else {                            // 1    ; 6
                System.out.println("m is very not");
            }
        } else {                    // 1    ; 7
            System.out.println("m is really not");
        }

        if (m > 0 && m % 2 == 0 && boolMethod()) {  // 8
            System.out.println("m is");
        }
    } // total 8



    public void elFor1(int j) {
        for (int i = 0; i < j; i++) {
            System.out.println("HEY");
        }
    }


    public void elFor2(int j) {
        for (int i = 0; i < j; i++) {
            for (int k = 0; k < j && j > 100; k++) {
                System.out.println("HEY");
            }
        }
    }


    public void elFor3(int j) {
        for (int i = 0; i < j; i++)
            for (int k = 0; k < j && j > 100; k++)
                System.out.println("HEY");
    }


    public void elFor4(int j) {                        // Cyc ;1      Cog ;1
        for (int i = 0; i < j; i++)                    // 1   ;2      1   ;2
            for (int k = 0; k < j && j > 100; k++)     // 2   ;4      2   ;4
                if (k != 10){                          // 1   ;5      3   ;7
                    System.out.println("HEY");         //
                    if (k != 9)                        // 1   ;6      4   ;11
                        System.out.println("HEY");     //
                    else                               //             1   ;12
                        System.out.println("HEY no");  //
                }                                      //
    }


    public void elWhile1(int w) {
        while (w > 0) {
            w--;
        }
    } // cyc and cog is 2


    public void elWhile2(int w) {
        while (w > 0)
            w--;

    }


    public void elWhile3(int w) {               // cyc; 1 cog ;1
        while (w > -100) {                      //  1;2       1;2
            boolean b = w > -1 && w < -10;      //  1;3
            while (b) {                         //  1;4       2;4
                w += (-1 * w * 80) / w *2;
            }
        }
    }


    public void elWhile4(int w) {               // cyc; 4 cog ;4
        while (w > -100)
            while (w > -1 && w < -10)
                w += (-1 * w * 80) / w *2;
    }


    public void elWhile5(int w) {               // cyc; 4 cog ;7
        while (w > -100)
            while ( w > -1 )
                while ( w > -10)
                    w += (-1 * w * 80) / w *2;
    }

    public void elDoWhi1(int i, int j) { // 1;1  1;1
        do {                             // 2, 2
            i++;
        } while (i < j);
    }

    public void elDoWhi2(int i, int j) { // 2  2
        do i++; while (i < j);
    }

    public void elTryCatch1(int i) { // 1           1
        if (i > -100)               // 1;2          1;2
            try {                       //          0; 2
                i++;                    //
            } catch (Exception ex) {    //  1;3    2;4
                if (i > 0) {            //  1;4    3;7
                    System.out.println();
                }
            }

    }


    public void elif1(int x) {                  // Cyc;1        Cog;1
        if (x == 1) {                           // 1;2        1;2
            System.out.println("hey");
        } else {                                // 0;2        1;3
            if (x ==2 ){                        // 1;3        2;5
                System.out.println("hey");
            } else if (x > 10){                 // 1;4        1;6
                while (x == 100)                // 1;5        3;9
                    System.out.println("hey");
            }
            else {                              // 0;5        1;10
                System.out.println("hey");
            }
        }
    }

    public void switch1(int i) {
        switch (i) {                                // 1     2
            case 1 : System.out.println("Hey");    // 2
                     break;
            case 2 : System.out.println("Man");    //3
                    break;
            default: System.out.println("Dude");    //4
        }
    }

    public String toString1(String[] array){
        String output = "<";
        boolean isNextOccupied = false;
        for (int i = 0; i < array.length && isNextOccupied == false; i++){
            output += array[i];
            if((i + 1) != array.length) {
                if (array[i + 1] == null) {
                    isNextOccupied = true;
                } else {
                    output += ", ";
                }
            }
        }
        output += ">";
        return output;
    }


    public static void horribleMethod(int i) {      // Cyc;1    Cog;1
        if (i == 10) {                              // 1;2      1;2
            class m {
                int i;
                m() {
                   i = 10;
                }

                public boolean iTime() {
                    if (i > 10)  {
                        return i + 1 == 15;
                    }
                    else
                        return i < 1;
                }
            }

            m m = new m();
            if (m.iTime())                       // 1;3       2;4
                System.out.println("wow");
            else {                               // 0;3       1;5
                if (i != m.i)                    // 1;4       3;8
                    System.out.println("yipee");
            }

        }

    }

    void doWhile(int z) {           //      COG: 1    CYC: 1
        do                          //      1;2         1;2
            System.out.println();
        while (z-- > 0);            // 0,0


        do {                        //      1;3         1;3
            System.out.println();
        }



        while (z-- > 0 && z <= -100);//     0;3         1;4


        while (z != 1)               //     1;4         1;5
            do                       //     2;6         1;6
                if (z++ == 20)       //     3;9         1;7
                    System.out.println("hey");
                else                 //     1;10           0;7
                    System.out.println("hey");
            while (z < 100  && z == -100);   // 0;            1;8

    }


    void doWhile2(int m) {
        do ; while (m > 1);
    }

    void doWhile3(int m) {
        do
            do
                ;
            while (m > 1);
        while (m > 1);
    }

    void doWhile4(int m) {
        do
            do
                if (m == 4)
                    System.out.println();
                else
                    System.out.println();
            while (m > 1);
        while (m > 1);
    }



    static boolean boolMethod() {
        return true;
    }


    void not() {}
    String str() { return "hey there joe"; }

    int methodName() { return 1; }

    public static void loop(int i) {
        for (int j = 0; j < i; j++) {
            System.out.println("loop");
        }
    }

    String typesofif(int n, String s) { // 1;1
        if (1==n) { // 2;2

        }

        while (n != 11) { // 3;3

        }

        for (int i = 0; i < n; i++) { // 4;4

        }

        do {        // 5;5

        } while(n > 1);

        try {

        } catch (Exception ex) { // 6;6

        }

        return switch (s) {      //6;7
            //case null -> "n";
            case "a" -> "";      //7;7
            case "b", "c" -> "a"; //8;7
            default -> "o";      //9;7
        };

    }

    void recursiveMethod(int n) {
        if (n > 0) {
            return;
        }
        recursiveMethod(n-1);
    }

    void elifTime(int m) {
        if (m == 1) {

        }
        else if (m == 2) {

        }
        else {

        }
    }

    class p {
        public static class k {
            void k_method(int m) {
                if (m < 0) {

                }

                class l {
                    public int n;

                    l (int n ) {
                        this.n = n;
                    }
                    int add(int i) {
                        if (i != 0) {
                            return 1;
                        } else if (i > 20) {
                            return 0;
                        } else {
                            i = i + 1;
                        }
                        return i + n;
                    }
                }

                l x = new l(1);

                if (x.n > 1000) {
                    x.add(-1000);
                }

            }
        }
    }

    class node {
        int val;
        node d;

        node() {
            val = 1;
            d = null;
        }

        void setVal(int n) {
            if (n < 0) {
                n = n * -1;
            }
            val = n;
        }

    }

}


abstract class abstractClass {
    abstract void thisShouldntBeIncluded();

    int thisShouldBeIncluded(int m) {
        if (m >1)
            if (m > 2)
                if (m> 3) {
                    return 2;
                }
        return 1;
    }

}


interface interfaceClass {
    default int isIncluded() {
        return 1;
    }

    int notIncluded();

}


class extra {

    private String[] array;
    private static final int CAPACITY = 10;
    private int size;


    // This should be cog = 2, but ended up being cog  = 8. That's because the recursive check did not
    // check if what was being called was a method or a variable, so a variable with the same name
    // as the method makes it think a recursive call is being done
    public String toString() {
        String toString = "<";
        for(int i = 0; i<size-1;i++){
            toString = toString.concat(array[i] + ", ");
        }

        toString = toString.concat(array[size-1] + ">");
        return toString;
    }
}

